package io.clawdroid.setup

import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color
import io.clawdroid.core.ui.theme.GlassBorder
import io.clawdroid.core.ui.theme.GlassWhite
import io.clawdroid.core.ui.theme.NeonCyan
import io.clawdroid.core.ui.theme.TextPrimary

@Composable
fun setupFieldColors() = OutlinedTextFieldDefaults.colors(
    focusedBorderColor = NeonCyan.copy(alpha = 0.5f),
    unfocusedBorderColor = GlassBorder,
    focusedContainerColor = GlassWhite,
    unfocusedContainerColor = Color.Transparent,
    focusedTextColor = TextPrimary,
    unfocusedTextColor = TextPrimary,
)
