import RegisterForm from "@/components/auth/register-form";
import { useAuth } from "@/context/auth-context";

const Register = () => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isAuthenticated && !isLoading) {
        return null;
    }

    return (
        <RegisterForm />
    );
};

export default Register;
