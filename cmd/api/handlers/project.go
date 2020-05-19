package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSameContracts godoc
// @Summary Get same contracts
// @Description Get same contracts
// @Tags contract
// @ID get-contract-same
// @Param network path string true "Network"
// @Param address path string true "KT address"
// @Param offset query integer false "Offset"
// @Param size query integer false "Requested count"
// @Accept json
// @Produce json
// @Success 200 {object} SameContractsResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /contract/{network}/{address}/same [get]
func (ctx *Context) GetSameContracts(c *gin.Context) {
	var req getContractRequest
	if err := c.BindUri(&req); handleError(c, err, http.StatusBadRequest) {
		return
	}

	var pageReq pageableRequest
	if err := c.BindQuery(&pageReq); handleError(c, err, http.StatusBadRequest) {
		return
	}

	by := map[string]interface{}{
		"address": req.Address,
		"network": req.Network,
	}
	contract, err := ctx.ES.GetContract(by)
	if handleError(c, err, 0) {
		return
	}

	sameContracts, err := ctx.ES.GetSameContracts(contract, 0, pageReq.Offset)
	if handleError(c, err, 0) {
		return
	}

	var response SameContractsResponse
	response.FromModel(sameContracts)

	c.JSON(http.StatusOK, response)
}

// GetSimilarContracts godoc
// @Summary Get similar contracts
// @Description Get similar contracts
// @Tags contract
// @ID get-contract-similar
// @Param network path string true "Network"
// @Param address path string true "KT address"
// @Accept  json
// @Produce  json
// @Success 200 {array} SimilarContract
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /contract/{network}/{address}/similar [get]
func (ctx *Context) GetSimilarContracts(c *gin.Context) {
	var req getContractRequest
	if err := c.BindUri(&req); handleError(c, err, http.StatusBadRequest) {
		return
	}

	by := map[string]interface{}{
		"address": req.Address,
		"network": req.Network,
	}
	contract, err := ctx.ES.GetContract(by)
	if handleError(c, err, 0) {
		return
	}

	similar, err := ctx.ES.GetSimilarContracts(contract)
	if handleError(c, err, 0) {
		return
	}

	response := make([]SimilarContract, len(similar))
	for i := range similar {
		diff, err := ctx.getContractCodeDiff(
			CodeDiffLeg{Address: contract.Address, Network: contract.Network},
			CodeDiffLeg{Address: similar[i].Address, Network: similar[i].Network},
		)
		if handleError(c, err, 0) {
			return
		}
		response[i].FromModels(similar[i], diff)
	}

	c.JSON(http.StatusOK, response)
}

// GetProjects godoc
// @Summary Show projects
// @Description Get all projects stats
// @Tags projects
// @ID get-projects
// @Accept  json
// @Produce  json
// @Success 200 {array} ProjectStats
// @Failure 500 {object} Error
// @Router /projects [get]
func (ctx *Context) GetProjects(c *gin.Context) {
	projects, err := ctx.ES.GetProjectsStats()
	if handleError(c, err, 0) {
		return
	}
	stats := make([]ProjectStats, len(projects))
	for i := range projects {
		stats[i].FromModel(projects[i])
	}
	c.JSON(http.StatusOK, stats)
}
