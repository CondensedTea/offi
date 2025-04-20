export class APIError extends Error {
    public status: number
    public traceparent: string | null

    constructor(message: string, status: number, traceparent: string | null) {
        super(message);

        this.status = status
        this.traceparent = traceparent
    }

    toString(): string {
        return `[${this.status}] API error: ${this.message}\nTrace: ${this.traceparent})`
    }

    static async fromResponse(r: Response): Promise<APIError> {
        return new APIError(await r.text(), r.status, r.headers.get("traceparent"));
    }
}

export const NoRecruitmentInfo = new Error("this team doesn't have recruitment post");
export const MatchNotFound = new Error("match not found or was not played yet");
export const NoLogsError = new Error("api didn't found logs for this match");
