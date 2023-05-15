<script lang="ts">
	import icon from "$lib/assets/icon.png";
	import telegram_logo from "$lib/assets/telegram-logo.svg";
	import instagram from "$lib/assets/instagram.svg";
	import facebook from "$lib/assets/facebook.svg";
	import arrow_down from "$lib/assets/arrow-down.svg";
	import header_bg from "$lib/assets/header_bg.jpg";
	import NavbarLinks from "$lib/components/NavbarLinks.svelte";
	import Page from "$lib/components/Page.svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { onMount } from "svelte";

	function smooth_scroll(event) {
		event.preventDefault()
		const link = event.currentTarget
		const anchorId = new URL(link.href).hash.replace('#', '')
		const anchor = document.getElementById(anchorId)
		window.scrollTo({
			top: anchor.offsetTop,
			behavior: 'smooth'
		})
	}

	onMount(async () => {
		/* particlesJS.load(@dom-id, @path-json, @callback (optional)); */
		// is loaded as a static script in app.html
		// @ts-ignore
		particlesJS("particles-js", {
			particles: {
				number: {
					value: 30,
					density: {
						enable: true,
						value_area: 800,
					},
				},
				color: {
					value: "#ffffff",
				},
				shape: {
					type: "polygon",
					stroke: {
						width: 0,
						color: "#000000",
					},
					polygon: {
						nb_sides: 5,
					},
					image: {
						src: "img/github.svg",
						width: 100,
						height: 100,
					},
				},
				opacity: {
					value: 0.5,
					random: false,
					anim: {
						enable: false,
						speed: 1,
						opacity_min: 0.1,
						sync: false,
					},
				},
				size: {
					value: 3,
					random: true,
					anim: {
						enable: false,
						speed: 19.18081918081918,
						size_min: 0.1,
						sync: false,
					},
				},
				line_linked: {
					enable: true,
					distance: 150,
					color: "#ffffff",
					opacity: 0.4,
					width: 1,
				},
				move: {
					enable: true,
					speed: 4,
					direction: "none",
					random: true,
					straight: false,
					out_mode: "out",
					bounce: false,
					attract: {
						enable: false,
						rotateX: 600,
						rotateY: 1200,
					},
				},
				nb: 80,
			},
			interactivity: {
				detect_on: "canvas",
				events: {
					onhover: {
						enable: false,
						mode: "grab",
					},
					onclick: {
						enable: true,
						mode: "push",
					},
					resize: true,
				},
				modes: {
					grab: {
						distance: 400,
						line_linked: {
							opacity: 1,
						},
					},
					bubble: {
						distance: 400,
						size: 40,
						duration: 2,
						opacity: 8,
						speed: 3,
					},
					repulse: {
						distance: 200,
						duration: 0.4,
					},
					push: {
						particles_nb: 4,
					},
					remove: {
						particles_nb: 2,
					},
				},
			},
			retina_detect: true,
		});
	});
</script>

<div id="particles-js" style={`background-image: url(${header_bg});`}>
	<div class="header">
		<img
			src={icon}
			alt={l($lang, $iso, "ui_logo_alt")}
			width="150"
			style="margin: auto; filter: invert(1);"
		/>
		<h1 style="margin-top:0; padding-top:0;">
			<span class="site-title">
				{l($lang, $iso, "ui_title")}
			</span>
			<span class="site-description">
				{l($lang, $iso, "ui_subtitle")}
			</span>
		</h1>
		<div class="header-icons">
			<a
				href="https://instagram.com/hay.rave"
			>
				<img
					src={instagram}
					alt={l($lang, $iso, "ui_instagram_logo_alt")}
					class="mt-2 mr-2"
					width="30"
				/>
			</a>
			<a
				href="https://t.me/+1IPjlBpQdzRmYjMy"
			>
				<img
					src={telegram_logo}
					alt={l($lang, $iso, "ui_telegram_logo_alt")}
					class="mt-2"
					width="30"
				/>
			</a>
			<a
				href="https://www.facebook.com/hayravefestival"
			>
				<img
					src={facebook}
					alt={l($lang, $iso, "ui_facebook_logo_alt")}
					class="mt-2 ml-2"
					width="30"
				/>
			</a>
		</div>
		<div class="header-links">
			<NavbarLinks />
		</div>
	</div>
	<a class="down" href="#topic1" on:click={smooth_scroll}>
		<img
			src={arrow_down}
			style="filter: invert(1); margin: auto;"
			alt={l($lang, $iso, "ui_arrow_down_alt")}
			class="mt-2"
			width="30"
		/>
	</a>
</div>
<div id="topic1" class="content">
	<Page path="homepage" />
</div>

<style>
	#particles-js {
		background-position: center;
		background-repeat: no-repeat;
		background-size: cover;
		display: flex;
		vertical-align: bottom;
		width: 100%;
		min-height: 100vh;
		background-repeat: no-repeat;
		background-size: cover;
		background-attachment: fixed;
	}
	#particles-js.particles-js-canvas-el {
		width: 100% !important;
		height: 100vh !important;
	}

	.header {
		z-index: 1;
		text-align: center;
		color: #fff;
		position: absolute;
		top: 50%;
		left: 50%;
		margin: 0 auto;
		-webkit-transform: translate(-50%, -50%);
		transform: translate(-50%, -50%);
	}
	.site-title {
		font-size: 50px;
		display: block;
		line-height: 1;
		color: #fff;
	}
	.site-description {
		font-size: 20px;
		display: block;
		line-height: 1;
		color: #fff;
		margin-top: 10px;
	}

	.header-links {
		margin: 10px;
	}
	.header-links.link {
		color: #fff;
		text-decoration: none;
		font-size: 15px;
		margin: 10px;
	}

	.header-icons {
		display: flex;
		justify-content: center;
		text-align: center;
	}
	.header-icons.icon {
		text-align: center;
		color: #fff;
		width: 20px;
		height: 20px;
		font-size: 20px;
		padding: 10px;
		margin: 5px;
		border-radius: 50%;
		border: 2px solid #fff;
		transition: all 0.7s;
	}
	.header-icons.icon:hover {
		color: #1a222c;
		background: #fff;
	}
	.header-icons.icon:active {
		color: #1a222c;
		background: #fff;
	}

	.down {
		color: #fff;
		position: absolute;
		bottom: 25px;
		width: 100%;
		margin: 0 auto;
		display: block;
		text-align: center;
		font-size: 30px;
		cursor: pointer;
	}
	.down.icon {
		position: absolute;
		-webkit-transform-style: preserve-3d;
		transform-style: preserve-3d;
		top: 50%;
		left: 50%;
		-webkit-transform: translate(-50%, -50%);
		transform: translate(-50%, -50%);
		width: 100px;
		height: 100px;
		fill: #fff;
		-webkit-animation: pulse 1.3s infinite;
		animation: pulse 1.3s infinite;
	}
	.down.icon:hover {
		color: #4b5664;
	}
	.down.icon:active {
		color: #4b5664;
	}
</style>
