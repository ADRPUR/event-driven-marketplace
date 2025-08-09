import dayjs from "dayjs";

export interface User {
  id: string;
  email: string;
  role: string;
  firstName?: string;
  lastName?: string;
  dateOfBirth?: string | dayjs.Dayjs;
  phone?: string;
  address?: Address | null;
  photo?: string;
  thumbnail?: string;
}

export interface Address {
  line?: string;
  city?: string;
  postal_code?: string;
  country?: string;
}

export interface UserDetails {
  id: string;
  email: string;
  role: string;
  details: {
    FirstName?: string;
    LastName?: string;
    DateOfBirth?: string | dayjs.Dayjs;
    Phone?: string;
    Address?: Address | null;
    PhotoPath?: string;
    Thumbnail?: string;
    CreatedAt?: string | dayjs.Dayjs;
    UpdatedAt?: string | dayjs.Dayjs;
  }
}

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  sessionToken: string;
  expiresAt: string | number;
  user: User;
}

export interface UsersListResponse {
  users: User[];
}

export interface UpdateProfileRequest {
  firstName?: string;
  lastName?: string;
  dateOfBirth?: string;
  phone?: string;
  address?: string;
  photoPath?: string;
  thumbnailPath?: string;
}


export interface UploadPhotoResponse {
  photoPath: string;
  thumbnailPath?: string;
}
