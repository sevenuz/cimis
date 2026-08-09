import type { RecordModel } from "pocketbase";

export interface User extends RecordModel {
	username: string;
	admin: boolean;
}
