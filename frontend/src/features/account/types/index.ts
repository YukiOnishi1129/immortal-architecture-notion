export interface Account {
  id: string;
  firstName: string;
  lastName: string;
  fullName: string;
  email: string;
  thumbnail: string | null;
  lastLoginAt: string | null;
  createdAt: string;
  updatedAt: string;
}
