// packages/frontend/validators/auth.ts
import { z } from 'zod';

export const LoginSchema = z.object({
    username: z.string().min(3, "Username minimal terdiri dari 3 karakter"),
    password: z.string().min(6, "Password minimal terdiri dari 6 karakter")
});

export type LoginFormData = z.infer<typeof LoginSchema>;