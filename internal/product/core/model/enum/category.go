package enum

type CategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Category uint

const (
	Lanche Category = iota
	Acompanhamento
	Bebida
	Sobremesa
)

var categoryToString = map[Category]string{
	Lanche:         "lanche",
	Acompanhamento: "acompanhamento",
	Bebida:         "bebida",
	Sobremesa:      "sobremesa",
}

var stringToCategory = map[string]Category{
	"lanche":         Lanche,
	"acompanhamento": Acompanhamento,
	"bebida":         Bebida,
	"sobremesa":      Sobremesa,
}

func (c Category) String() string {
	if name, ok := categoryToString[c]; ok {
		return name
	}
	return "unknown"
}

func IsValidCategory(value uint) bool {
	_, ok := categoryToString[Category(value)]
	return ok
}

func FromCategoryString(name string) (Category, bool) {
	cat, ok := stringToCategory[name]
	return cat, ok
}

func GetAllCategories() []CategoryDTO {
	categories := make([]CategoryDTO, 0, len(categoryToString))
	for cat, name := range categoryToString {
		categories = append(categories, CategoryDTO{
			ID:   uint(cat),
			Name: name,
		})
	}
	return categories
}
