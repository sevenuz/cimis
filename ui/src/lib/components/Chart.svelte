<script lang="ts">
	import { onMount, onDestroy } from "svelte";
	import {
		Chart as ChartJS,
		BarController,
		LineController,
		BarElement,
		LineElement,
		PointElement,
		LinearScale,
		CategoryScale,
		Tooltip,
		Legend,
		type ChartData,
		type ChartOptions,
		type ChartType,
	} from "chart.js";

	ChartJS.register(
		BarController,
		LineController,
		BarElement,
		LineElement,
		PointElement,
		LinearScale,
		CategoryScale,
		Tooltip,
		Legend
	);

	export let type: ChartType;
	export let data: ChartData;
	export let options: ChartOptions = {};
	export let height = 260;

	let canvas: HTMLCanvasElement;
	let chart: ChartJS | null = null;

	onMount(() => {
		chart = new ChartJS(canvas, { type, data, options });
	});

	onDestroy(() => {
		chart?.destroy();
	});

	$: if (chart) {
		chart.data = data;
		chart.options = options;
		chart.update();
	}
</script>

<div style="height: {height}px;">
	<canvas bind:this={canvas} />
</div>
