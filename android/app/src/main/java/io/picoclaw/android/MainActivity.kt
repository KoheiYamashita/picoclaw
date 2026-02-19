package io.picoclaw.android

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import io.picoclaw.android.core.ui.theme.PicoClawTheme
import io.picoclaw.android.feature.chat.screen.ChatScreen
import io.picoclaw.android.feature.chat.screen.SettingsScreen
import io.picoclaw.android.navigation.NavRoutes

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            PicoClawTheme {
                val navController = rememberNavController()
                NavHost(navController = navController, startDestination = NavRoutes.CHAT) {
                    composable(NavRoutes.CHAT) {
                        ChatScreen(
                            onNavigateToSettings = { navController.navigate(NavRoutes.SETTINGS) }
                        )
                    }
                    composable(NavRoutes.SETTINGS) {
                        SettingsScreen(
                            onNavigateBack = { navController.popBackStack() }
                        )
                    }
                }
            }
        }
    }
}
