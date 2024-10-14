package def

const (
	ProductImagePackaging   = iota + 1 // 外部包装图
	ProductImageInternal               // 内部图
	ProductImageBarcode                // 条形码图
	ProductImageIngredients            // 配料表图
	ProductImageNutrition              // 营养成分图
	ProductImageOther                  // 其他图
)

var ProductImageTypes = map[int]string{
	ProductImagePackaging:   "外部包装图",
	ProductImageInternal:    "内部图",
	ProductImageBarcode:     "条形码图",
	ProductImageIngredients: "配料表图",
	ProductImageNutrition:   "营养成分图",
	ProductImageOther:       "其他图",
}
