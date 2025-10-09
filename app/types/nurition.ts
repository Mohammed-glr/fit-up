interface FoodLogEntry {
  id: string;
  userId: string;
  date: string;
  mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  foodItem: string;
  quantity: number;
  calories: number;
  protein: number;
  carbs: number;
  fat: number; 
}

type ListFoodLogEntriesResponse = {
  entries: FoodLogEntry[];
};

type CreateFoodLogEntryRequest = {
  date: string;
  mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  foodItem: string;
  quantity: number;
  calories: number;
  protein: number;
  carbs: number;
  fat: number; 
};

interface Nutrition {
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
}

