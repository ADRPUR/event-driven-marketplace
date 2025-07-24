#!/bin/bash
cd src

mkdir -p api assets components/common components/layout features/auth features/profile features/products features/users hooks store routes types

touch api/auth.ts api/user.ts api/product.ts
touch components/common/Loader.tsx
touch components/layout/Header.tsx components/layout/Footer.tsx components/layout/Sidebar.tsx components/layout/ProtectedRoute.tsx
touch features/auth/LoginPage.tsx features/auth/RegisterPage.tsx features/auth/useLogin.ts features/auth/useRegister.ts
touch features/profile/ProfilePage.tsx features/profile/EditProfileForm.tsx features/profile/useProfile.ts
touch features/products/ProductsPage.tsx
touch features/users/UsersTable.tsx features/users/EditUserDialog.tsx features/users/useUsers.ts
touch routes/AppRouter.tsx
touch hooks/useAuth.ts hooks/useFetch.ts
touch store/authStore.ts store/userStore.ts
touch types/user.ts types/auth.ts types/product.ts

echo "console.log('Hello from Vite!');" > assets/avatar_placeholder.png
echo "export {};" > api/auth.ts

