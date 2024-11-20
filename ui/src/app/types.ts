export type BlogPost = {
    id: number;
    title: string;
    content: string;
    status: string;
    created_at: string;
    account_id: number;
    author: string;
};

export type User = {
    exp: number
    id: number
    role: number // User: 0, Admin: 1
    username: string
    authToken: string
}

export type UserEntry = {
    id: number
    username: string
    role: number
}

export type statusResponseResult = {
    initialized: boolean
}
