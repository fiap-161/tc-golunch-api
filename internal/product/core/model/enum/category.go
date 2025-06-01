package enum

type CategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Category uint

const (
	Unknown Category = iota
	Meal
	Side
	Drink
	Dessert
)

var categoryToString = map[Category]string{
	Unknown: "unknown",
	Meal:    "meal",
	Side:    "side",
	Drink:   "drink",
	Dessert: "dessert",
}

var stringToCategory = map[string]Category{
	"unknown": Unknown,
	"meal":    Meal,
	"side":    Side,
	"drink":   Drink,
	"dessert": Dessert,
}

func (c Category) String() string {
	if name, ok := categoryToString[c]; ok {
		return name
	}
	return "unknown"
}

func IsValidCategory(value uint) bool {
	cat := Category(value)
	_, ok := categoryToString[cat]
	return ok && cat != Unknown
}

func FromCategoryString(name string) (Category, bool) {
	cat, ok := stringToCategory[name]
	return cat, ok
}

func GetAllCategories() []CategoryDTO {
	categories := make([]CategoryDTO, 0, len(categoryToString))
	for cat, name := range categoryToString {
		if cat == Unknown {
			continue
		}
		categories = append(categories, CategoryDTO{
			ID:   uint(cat),
			Name: name,
		})
	}
	return categories
}
