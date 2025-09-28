import * as SecureStore from "expo-secure-store";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { Platform } from "react-native";

const isSecureStorageAvailable = Platform.OS !== "web";

class SecureStorage {
    async setToken(key: string, value: string): Promise<void> {
        if (isSecureStorageAvailable) {
            await SecureStore.setItemAsync(key, value);
        } else {
            await AsyncStorage.setItem(key, value);
        }
    }

    async getToken(key: string): Promise<string | null> {
        if (isSecureStorageAvailable) {
            return await SecureStore.getItemAsync(key);
        } else {
            return await AsyncStorage.getItem(key);
        }
    }

    async removeToken(key: string): Promise<void> {
        if (isSecureStorageAvailable) {
            await SecureStore.deleteItemAsync(key);
        } else {
            await AsyncStorage.removeItem(key);
        }
    }

    async clearTokens(): Promise<void> {
        await Promise.all([
            this.removeToken('access_token'),
            this.removeToken('refresh_token')
        ]);
    } 
}

export const secureStorage = new SecureStorage();
