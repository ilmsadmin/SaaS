#!/bin/bash

echo "📱 Setting up Zplus Mobile App with React Native"
echo "==============================================="

# Navigate to mobile directory
cd /Users/toan/Documents/project/SaaS/apps/frontend/mobile

# Initialize React Native with Expo
npx create-expo-app@latest . --template blank-typescript

# Install additional dependencies
npm install @tanstack/react-query axios react-navigation/native react-navigation/native-stack @expo/vector-icons react-native-safe-area-context react-native-screens react-native-gesture-handler @react-native-async-storage/async-storage react-hook-form react-native-toast-message

# Install Expo specific packages
npx expo install expo-router expo-constants expo-status-bar

echo "✅ Mobile app setup completed!"
echo "📱 Mobile app is ready at: apps/frontend/mobile"
echo "🚀 Run 'cd apps/frontend/mobile && npm start' to start development"
