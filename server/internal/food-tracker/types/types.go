package types

type RecipeCategory string

const (
	CategoryBreakfast RecipeCategory = "breakfast"
	CategoryLunch     RecipeCategory = "lunch"
	CategoryDinner    RecipeCategory = "dinner"
	CategorySnack     RecipeCategory = "snack"
	CategoryDessert   RecipeCategory = "dessert"
)

type RecipeDifficulty string

const (
	DifficultyEasy   RecipeDifficulty = "easy"
	DifficultyMedium RecipeDifficulty = "medium"
	DifficultyHard   RecipeDifficulty = "hard"
)

type SystemRecipe struct {
	RecipeID          int              `json:"id"`
	RecipeName        string           `json:"name"`
	RecipeDesc        string           `json:"description"`
	RecipesCategory   RecipeCategory   `json:"category"`
	RecipesDifficulty RecipeDifficulty `json:"difficulty"`
	RecipesCalories   int              `json:"calories"`
	RecipesProtein    int              `json:"protein"`
	RecipesCarbs      int              `json:"carbs"`
	RecipesFat        int              `json:"fat"`
	RecipesFiber      int              `json:"fiber"`
	PrepTime          int              `json:"prep_time"`
	CookTime          int              `json:"cook_time"`
	Servings          int              `json:"servings"`
	RecipesImageURL   string           `json:"image_url"`
	IsActive          bool             `json:"is_active"`
	CreatedAt         string           `json:"created_at"`
	UpdatedAt         string           `json:"updated_at"`
}

type SystemRecipesIngredient struct {
	IngredientID     int     `json:"id"`
	RecipeID         int     `json:"recipe_id"`
	IngredientItem   string  `json:"item"`
	IngredientAmount float64 `json:"amount"`
	IngredientUnit   string  `json:"unit"`
	OrderIndex       int     `json:"order_index"`
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
	RecipeID          int              `json:"id"`
	UserID            string           `json:"user_id"`
	RecipeName        string           `json:"name"`
	RecipeDesc        string           `json:"description"`
	RecipesCategory   RecipeCategory   `json:"category"`
	RecipesDifficulty RecipeDifficulty `json:"difficulty"`
	RecipesCalories   int              `json:"calories"`
	RecipesProtein    int              `json:"protein"`
	RecipesCarbs      int              `json:"carbs"`
	RecipesFat        int              `json:"fat"`
	RecipesFiber      int              `json:"fiber"`
	PrepTime          int              `json:"prep_time"`
	CookTime          int              `json:"cook_time"`
	Servings          int              `json:"servings"`
	RecipesImageURL   string           `json:"image_url"`
	IsFavorite        bool             `json:"is_favorite"`
	CreatedAt         string           `json:"created_at"`
	UpdatedAt         string           `json:"updated_at"`
}

type UserRecipesIngredient struct {
	IngredientID     int     `json:"id"`
	RecipeID         int     `json:"recipe_id"`
	IngredientItem   string  `json:"item"`
	IngredientAmount float64 `json:"amount"`
	IngredientUnit   string  `json:"unit"`
	OrderIndex       int     `json:"order_index"`
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
	Id        int    `json:"id"`
	UserID    string `json:"user_id"`
	RecipeID  int    `json:"recipe_id"`
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
	EntryID        int      `json:"id"`
	UserID         string   `json:"user_id"`
	LogDate        string   `json:"log_date"`
	SystemRecipeID *int     `json:"system_recipe_id,omitempty"`
	UserRecipeID   *int     `json:"user_recipe_id,omitempty"`
	Calories       int      `json:"calories"`
	Protein        int      `json:"protein"`
	Carbs          int      `json:"carbs"`
	Fat            int      `json:"fat"`
	Fiber          int      `json:"fiber"`
	Servings       float64  `json:"servings"`
	MealType       MealType `json:"meal_type"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type SystemRecipeDetail struct {
	SystemRecipe
	Ingredients  []SystemRecipesIngredient  `json:"ingredients"`
	Instructions []SystemRecipesInstruction `json:"instructions"`
	Tags         []SystemRecipesTag         `json:"tags"`
}

type UserRecipeDetail struct {
	UserRecipe
	Ingredients  []UserRecipesIngredient  `json:"ingredients"`
	Instructions []UserRecipesInstruction `json:"instructions"`
	Tags         []UserRecipesTag         `json:"tags"`
}

type UserAllRecipesView struct {
	Source     string         `json:"source"`
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Category   RecipeCategory `json:"category"`
	Calories   int            `json:"calories"`
	Protein    int            `json:"protein"`
	Carbs      int            `json:"carbs"`
	Fat        int            `json:"fat"`
	Fiber      int            `json:"fiber"`
	PrepTime   int            `json:"prep_time"`
	Servings   int            `json:"servings"`
	ImageURL   string         `json:"image_url"`
	UserID     *string        `json:"user_id,omitempty"`
	IsFavorite bool           `json:"is_favorite"`
}

type DailyNutritionSummary struct {
	UserID        string `json:"user_id"`
	LogDate       string `json:"log_date"`
	TotalCalories int    `json:"total_calories"`
	TotalProtein  int    `json:"total_protein"`
	TotalCarbs    int    `json:"total_carbs"`
	TotalFat      int    `json:"total_fat"`
	TotalFiber    int    `json:"total_fiber"`
	TotalEntries  int    `json:"total_entries"`
}

type FoodLogEntryWithRecipe struct {
	FoodLogEntry
	RecipeName   string `json:"recipe_name,omitempty"`
	RecipeSource string `json:"recipe_source,omitempty"`
}

type CreateRecipeRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Category    RecipeCategory   `json:"category"`
	Difficulty  RecipeDifficulty `json:"difficulty"`
	Calories    int              `json:"calories"`
	Protein     int              `json:"protein"`
	Carbs       int              `json:"carbs"`
	Fat         int              `json:"fat"`
	Fiber       int              `json:"fiber"`
	PrepTime    int              `json:"prep_time"`
	CookTime    int              `json:"cook_time"`
	ImageURL    string           `json:"image_url"`
	Servings    int              `json:"servings"`
	Ingredients []struct {
		IngredientID int     `json:"ingredient_id"`
		Item        string  `json:"item"`
		Amount     float64 `json:"amount"`
		Unit       string  `json:"unit"`
		OrderIndex int     `json:"order_index"`
	} `json:"ingredients"`
	Instructions []struct {
		InstructionID int    `json:"instruction_id"`
		StepNumber    int    `json:"step_number"`
		Instruction   string `json:"instruction"`
		Text          string `json:"text"`
	} `json:"instructions"`
	Tags []string `json:"tags"`
}

type CreateFoodLogRequest struct {
	LogDate        string   `json:"log_date"`
	MealType       MealType `json:"meal_type"`
	SystemRecipeID *int     `json:"system_recipe_id,omitempty"`
	UserRecipeID   *int     `json:"user_recipe_id,omitempty"`
	Calories       int      `json:"calories"`
	Protein        int      `json:"protein"`
	Carbs          int      `json:"carbs"`
	Fat            int      `json:"fat"`
	Fiber          int      `json:"fiber"`
	Servings       float64  `json:"servings"`
}

type RecipeFilters struct {
	Category    *RecipeCategory
	Difficulty  *RecipeDifficulty
	MaxCalories *int
	MinProtein  *int
	MaxPrepTime *int
	Tags        []string
	IsFavorite  *bool
	SearchTerm  string
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}

type SearchQuery struct {
	Term          string
	Category      *RecipeCategory
	Difficulty    *RecipeDifficulty
	MaxCalories   *int
	MinProtein    *int
	MaxPrepTime   *int
	Tags          []string
	IncludeSystem bool
	IncludeUser   bool
	FavoritesOnly bool
	Limit         int
	Offset        int
}
