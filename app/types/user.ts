interface User {
  id: string;
  username: string;
  name: string;
  bio: string;
  email: string;
  image: string;
  role: 'admin' | 'user';
  is_two_factor_enabled: boolean;
  created_at: string;
  updated_at: string;
}