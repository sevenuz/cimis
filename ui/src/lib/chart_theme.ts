import type { ChartOptions } from "chart.js";

// Dark-mode categorical palette (fixed hue order - never cycled/reassigned by
// rank) and chart chrome tokens, per the dataviz skill's reference palette.
// This app is dark-theme only (see app.postcss), so only the dark column is used.
export const categorical_palette = [
	"#3987e5", // blue
	"#d95926", // orange
	"#199e70", // aqua
	"#c98500", // yellow
	"#d55181", // magenta
	"#008300", // green
	"#9085e9", // violet
	"#e66767", // red
];

export const chart_ink = {
	text_secondary: "#c3c2b7",
	text_muted: "#898781",
	grid: "#2c2c2a",
	baseline: "#383835",
};

// Products/payment methods carry their own admin-set `color` field (reused
// elsewhere for buttons via get_colors) - identity follows the entity, so charts
// reuse that same color. Entities without one fall back to the fixed palette by
// stable index so repeat renders stay consistent.
export function entity_color(color: string | undefined, index: number): string {
	return color && color !== "" ? color : categorical_palette[index % categorical_palette.length];
}

export function base_chart_options(): ChartOptions {
	return {
		responsive: true,
		maintainAspectRatio: false,
		color: chart_ink.text_secondary,
		plugins: {
			legend: {
				labels: { color: chart_ink.text_secondary },
			},
			tooltip: {
				backgroundColor: "#0d0d0d",
				titleColor: "#ffffff",
				bodyColor: chart_ink.text_secondary,
				borderColor: chart_ink.baseline,
				borderWidth: 1,
			},
		},
		scales: {
			x: {
				ticks: { color: chart_ink.text_muted },
				grid: { color: chart_ink.grid },
				border: { color: chart_ink.baseline },
			},
			y: {
				ticks: { color: chart_ink.text_muted },
				grid: { color: chart_ink.grid },
				border: { color: chart_ink.baseline },
				beginAtZero: true,
			},
		},
	};
}
