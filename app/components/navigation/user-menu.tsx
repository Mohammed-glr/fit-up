import {
    View,
    Text, StyleSheet,
    Platform
} from "react-native";
import { MotiView, MotiText } from "moti";
import { COLORS } from "@/constants/theme";
import { useCurrentUser } from "@/hooks/user/use-current-user"; 
import React from "react";
import { Avatar } from "@/components/ui/avatar";
import { 
    DropdownMenu,
    DropdownTrigger,
    DropdownMenuItem
} from "@/components/ui/dropdown-menu";
import LogoutButton from "../auth/logout-button";

export const UserMenu: React.FC = () => {
    const { data: user, isLoading } = useCurrentUser();
    const [isOpen, setIsOpen] = React.useState(false);


    if (isLoading) return <div>Loading...</div>;
    if (!user) return <div>User not found</div>;
    return (
        <DropdownMenu isOpen={isOpen}>
            <DropdownTrigger
                label="User Menu"
                isOpen={isOpen}
                onPress={() => setIsOpen(!isOpen)}
            >
                <Avatar />
            </DropdownTrigger>
            <DropdownMenuItem
                item={{ label: user.name, value: user.id }}
                isSelected={false}
                onPress={() => {}}
                index={0}
            >
                <Text>{user.name}</Text>
            </DropdownMenuItem>
            <DropdownMenuItem
                item={{ label: "Logout", value: "logout" }}
                isSelected={false}
                onPress={() => {}}
                index={1}
            >
                <LogoutButton />
            </DropdownMenuItem>
        </DropdownMenu>
    );
}