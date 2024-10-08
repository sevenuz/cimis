import type { Record } from "pocketbase";
import type { PaymentMethod } from "./PaymentMethod";
import type { User } from "./User";

export interface Order extends Record {
	user: string; // relation
	payment_method: string; // relation
	total: number;
	is_bookout: boolean;
	expand: {
		user: User
		payment_method: PaymentMethod
	}
}
