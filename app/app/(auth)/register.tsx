import RegisterForm from "@/components/auth/register-form";
import { useCurrentUser } from "@/hooks/user/use-current-user";
import { Redirect } from "expo-router";

const Register = () => {
    const { data: currentUser, isLoading } = useCurrentUser();

    if (!isLoading && currentUser) {
        if (currentUser.role === 'coach') {
            return <Redirect href="/(coach)" />;
        }
        return <Redirect href="/(user)" />;
    }

    return (
        <RegisterForm />
    );
};

export default Register;
