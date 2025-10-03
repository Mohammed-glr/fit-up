package types



type RecipeCategory string

const (
	CategoryBreakfast  RecipeCategory = "breakfast"
	CategoryLunch      RecipeCategory = "lunch"
	CategoryDinner     RecipeCategory = "dinner"
	CategorySnack      RecipeCategory = "snack"
	CategoryDessert    RecipeCategory = "dessert"
)

type RecipeDifficulty string

const (
	DifficultyEasy   RecipeDifficulty = "easy"
	DifficultyMedium RecipeDifficulty = "medium"
	DifficultyHard   RecipeDifficulty = "hard"
)
type SystemRecipe struct {
	RecipeID   int    `json:"id"`
	RecipeName string `json:"name"`
	RecipeDesc string `json:"description"`
	RecipesCategory RecipeCategory `json:"category"`
	RecipesDifficulty RecipeDifficulty `json:"difficulty"`
	RecipesCalories int    `json:"calories"`
	RecipesProtein  int    `json:"protein"`
	RecipesCarbs    int    `json:"carbs"`
	RecipesFats     int    `json:"fats"`
	RecipesFiber   int    `json:"fiber"`
	PrepTimeMinutes int    `json:"prep_time_minutes"`
	CookTimeMinutes int    `json:"cook_time_minutes"`
	RecipesImageURL string `json:"image_url"`
	IsActive      bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}


type SystemRecipesIngredient struct {
	IngredientID   int    `json:"id"`
	RecipeID       int    `json:"recipe_id"`
	IngredientItem string `json:"item"`
	IngredientAmount float64 `json:"amount"`
	IngredientUnit string `json:"unit"`
	OrderIndex     int    `json:"order_index"`
}

type SystemRecipesInstruction struct {
	InstructionID   int    `json:"id"`
	RecipeID        int    `json:"recipe_id"`
	InstructionStep int    `json:"step_number"`
	InstructionText string `json:"instruction"`
}


type SystemRecipesTag struct {
	TagID    int    `json:"id"`
	RecipeID int    `json:"recipe_id"`
	TagName  string `json:"tag_name"`
}

type UserRecipe struct {
	RecipeID   int    `json:"id"`
	RecipeName string `json:"name"`
	RecipeDesc string `json:"description"`
	RecipesCategory RecipeCategory `json:"category"`
	RecipesDifficulty RecipeDifficulty `json:"difficulty"`
	RecipesCalories int    `json:"calories"`
	RecipesProtein  int    `json:"protein"`
	RecipesCarbs    int    `json:"carbs"`
	RecipesFats     int    `json:"fats"`
	RecipesFiber   int    `json:"fiber"`
	PrepTimeMinutes int    `json:"prep_time_minutes"`
	CookTimeMinutes int    `json:"cook_time_minutes"`
	RecipesImageURL string `json:"image_url"`
	IsFavorite      bool   `json:"is_favorite"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}


type UserRecipesIngredient struct {
	IngredientID   int    `json:"id"`
	RecipeID       int    `json:"recipe_id"`
	IngredientItem string `json:"item"`
	IngredientAmount float64 `json:"amount"`
	IngredientUnit string `json:"unit"`
	OrderIndex     int    `json:"order_index"`
}	


type UserRecipesInstruction struct {
	InstructionID   int    `json:"id"`
	RecipeID        int    `json:"recipe_id"`
	InstructionStep int    `json:"step_number"`
	InstructionText string `json:"instruction"`
}

type UserRecipesTag struct {
	TagID    int    `json:"id"`
	RecipeID int    `json:"recipe_id"`
	TagName  string `json:"tag_name"`
}


type UserFavoriteRecipe struct {
	Id       int `json:"id"`
	UserID   int `json:"user_id"`
	RecipeID int `json:"recipe_id"`
	CreatedAt string `json:"created_at"`
}



type MealType string

const (
	MealTypeBreakfast MealType = "breakfast"
	MealTypeLunch     MealType = "lunch"
	MealTypeDinner    MealType = "dinner"
	MealTypeSnack     MealType = "snack"
)	
type FoodLogEntry struct {
	EntryID     int    `json:"id"`
	UserID      int    `json:"user_id"`
	LogDate    string `json:"log_date"`
	SystemRecipeID *int   `json:"system_recipe_id,omitempty"`
	UserRecipeID   *int   `json:"user_recipe_id,omitempty"`
	Calories      int    `json:"calories"`
	Protein      int    `json:"protein"`
	Carbs        int    `json:"carbs"`
	Fats         int    `json:"fats"`
	Fiber       int    `json:"fiber"`
	MealType    MealType `json:"meal_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

