import type { RecordModel } from "pocketbase";

export interface ProductStat extends RecordModel {
	event: string;
	bar: string;
	product: string;
	free: boolean;
	units: number;
	revenue: number;
}

export interface ProductUpsellStat extends RecordModel {
	event: string;
	product: string;
	order_count: number;
	avg_order_total: number;
}

export interface PaymentMethodStat extends RecordModel {
	event: string;
	bar: string;
	payment_method: string;
	order_count: number;
	volume: number;
	avg_order_value: number;
	raw_total: number;
	fee_cost: number;
}

export interface TimeBucketStat extends RecordModel {
	event: string;
	bar: string;
	bucket: string;
	order_count: number;
	revenue: number;
}

export interface ProductTimeStat extends RecordModel {
	event: string;
	product: string;
	bucket: string;
	units: number;
	revenue: number;
}
