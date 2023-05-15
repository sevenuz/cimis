import { writable } from "svelte/store";

export enum NotificationType {
	success = "#84C991",
	info = "#5bc0de",
	warning = "#f0ad4e",
	error = "#E26D69"
}

export interface Notification {
	id?: string;
	message: string;
	duration: number;
	type: NotificationType;
}

export const notifications = writable([] as Notification[]);

export function notify(n: Notification) {
	notifications.update(v => {
		const id = generate_id();
		v.push({
			...n,
			id
		});
		setTimeout(() => notifications.update(v2 => v2.filter(n2 => n2.id != id)), n.duration);
		return v;
	});
}

function generate_id(): string {
    return '_' + Math.random().toString(36);
};
