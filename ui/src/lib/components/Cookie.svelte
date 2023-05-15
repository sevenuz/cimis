<script lang="ts">
	import { l, lang, iso } from "$lib/stores/lang";
  import { onMount } from "svelte";
	import { fade } from "svelte/transition";

	const cookie_name = "cookiebanner";
	let show_banner = false;

	onMount(() => {
		show_banner = get_cookie(cookie_name) !== "true";
	});

	function click_banner() {
		show_banner = false;
		set_cookie(cookie_name, "true");
	}

	// cookie code from:
	// https://gist.github.com/joduplessis/7b3b4340353760e945f972a69e855d11
	function set_cookie(name: string, val: string) {
		const date = new Date();
		const value = val;

		// Set it expire in 30 days
		date.setTime(date.getTime() + 30 * 24 * 60 * 60 * 1000);

		// Set it
		document.cookie =
			name + "=" + value + "; expires=" + date.toUTCString() + "; path=/";
	}

	function get_cookie(name: string) {
		const value = "; " + document.cookie;
		const parts = value.split("; " + name + "=");

		if (parts.length == 2) {
			return parts.pop().split(";").shift();
		}
	}
</script>

{#if show_banner}
	<div class="cookies bg-green" transition:fade>
		<div class="content">
			{l($lang, $iso, "ui_cookie_banner")}
		</div>
		<button class="bg-white" on:click={click_banner}>
			{l($lang, $iso, "ui_cookie_banner_action")}
		</button>
	</div>
{/if}

<style>
	.cookies {
		position: fixed;
		bottom: 0;
		right: 0;
		left: 0;
		margin: 0 auto;
		padding: 0;
		z-index: 9999;
		display: flex;
		flex-direction: column;
		justify-content: flex-start;
		align-items: center;
	}
	.content {
		padding: 10px;
		display: block;
		color: white;
		font-weight: 500;
	}
</style>
