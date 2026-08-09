import type { RecordModel } from "pocketbase";
import type { PaymentMethod } from "./PaymentMethod";
import type { User } from "./User";
import type { Bar } from "./Bar";
import type { Event } from "./Event";

export interface Order extends RecordModel {
	user: string; // relation
	payment_method: string; // relation
	bar: string; // relation
	event: string; // relation
	total: number;
	is_bookout: boolean;
	client_id: string; // idempotency key, generated client-side
	expand: {
		user: User;
		payment_method: PaymentMethod;
		bar: Bar;
		event: Event;
	}
}
