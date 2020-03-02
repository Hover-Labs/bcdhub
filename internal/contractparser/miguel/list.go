package miguel

import (
	"github.com/aopoltorzhicky/bcdhub/internal/contractparser/consts"
	"github.com/aopoltorzhicky/bcdhub/internal/contractparser/meta"
	"github.com/tidwall/gjson"
)

type listDecoder struct{}

// Decode -
func (l *listDecoder) Decode(node gjson.Result, path string, nm *meta.NodeMetadata, metadata meta.Metadata) (interface{}, error) {
	res := make([]interface{}, 0)
	arr := node.Array()
	if len(arr) > 0 {
		subPath := "/l"
		if nm.Type == consts.SET {
			subPath = "/s"
		}
		for _, arg := range arr {
			data, err := michelineNodeToMiguel(arg, path+subPath, metadata)
			if err != nil {
				return nil, err
			}
			res = append(res, data)
		}
	}
	return res, nil
}