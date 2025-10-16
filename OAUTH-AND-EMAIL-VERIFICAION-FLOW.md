Part 1: OAuth with PKCE
Backend (Go)
Register your app with the OAuth provider:
Go to the developer console for your provider (e.g., Google Cloud Console) and create new OAuth credentials.
Set the "Application type" to iOS or Android.
Configure your redirect URIs. For mobile, this will be a custom URI scheme (e.g., com.yourapp.bundleid://oauth-callback).
You will receive a Client ID. Unlike web, you won't use a client secret in your mobile app.
Handle the callback endpoint:
Create a backend endpoint (e.g., /api/auth/callback) that accepts the code and code_verifier sent by your React Native app.
Use Go's golang.org/x/oauth2 package to exchange the authorization code and code verifier for an access token.
Validate the incoming request, exchange the tokens securely on your backend, and retrieve user information from the provider.
Upon success, generate and send your own secure JSON Web Token (JWT) back to the mobile app for your API authentication. 
Frontend (React Native)
Install react-native-app-auth:
This library handles the complex native-side PKCE flow.
sh
npm install react-native-app-auth
Wees voorzichtig met code.

Configure deep linking:
iOS: In Xcode, configure Associated Domains and add your custom URI scheme to your Info.plist. For testing with TestFlight, this setup is the same as a production app.
Android: Add an intent-filter for your custom URI scheme (com.yourapp.bundleid) in android/app/src/main/AndroidManifest.xml.
Implement the OAuth flow:
Use react-native-app-auth to initiate the flow. This will open the external browser or a provider-specific pop-up.
The user logs in on the provider's page. The provider then redirects to your app using the deep link.
react-native-app-auth intercepts the deep link, extracts the code and code_verifier, and performs the token exchange.
Pass the received code and code_verifier to your Go backend to get your own application's JWT.
Store the JWT securely (e.g., using a library like react-native-keychain). 
Part 2: Email verification
Backend (Go)
Generate a verification token:
When a new user signs up with an email, generate a unique, cryptographically secure token.
Associate this token with the user in your database and set an expiration time.
Send the email with Resend:
Use Resend's Go SDK to send an email to the user.
The email should contain a deep link using your app's custom URI scheme, embedding the token as a query parameter (e.g., com.yourapp.bundleid://verify?token=...).
Create the verification endpoint:
Create a protected API endpoint on your Go backend (e.g., /api/verify-email).
This endpoint should accept the verification token.
Look up the user associated with the token, check if the token is valid and not expired, and mark the user's email as verified in your database. 
Frontend (React Native)
Configure deep linking:
The deep linking setup is the same as for OAuth, using your custom URI scheme (com.yourapp.bundleid).
Use a navigation library like React Navigation to listen for incoming deep links.
Handle incoming links:
Configure your app's entry point to check for a deep link on startup.
If a link with a path like /verify and a token query parameter is detected:
Extract the token.
Navigate the user to a verification screen in your app.
Make an API call to your Go backend's verification endpoint, passing the token.
Show the user a confirmation message and navigate them to the main app screen. 
Final flow summary
OAuth:
User taps "Login with Google" in your React Native app.
react-native-app-auth initiates the PKCE flow, opening the browser for the user to log in.
The provider redirects the user back to your app via a deep link.
The React Native app's deep link handler catches the redirect and gets the authorization code.
The app sends the authorization code and PKCE verifier to your Go backend.
Your Go backend exchanges this for an access token, gets user data, and issues a JWT for your app.
Your Go backend sends your app's JWT back to the React Native app.
The React Native app stores and uses the JWT for future authenticated requests.
Email Verification:
User registers with an email address in your React Native app.
The app calls your Go backend's signup endpoint.
Your Go backend creates a new user, generates a unique verification token, and sends a deep link containing the token via Resend.
User clicks the link in their email, which opens your app.
The deep link is captured by your React Native app.
The app sends the token to your Go backend's verification endpoint.
The backend validates the token, marks the user as verified, and the frontend updates the UI accordingly. 