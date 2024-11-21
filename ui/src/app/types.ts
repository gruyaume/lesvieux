export type JobPost = {
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
    email: string
    authToken: string
}

export type UserEntry = {
    id: number
    email: string
    role: number
}

export type EmployerEntry = {
    id: number
    name: string
}


export type statusResponseResult = {
    initialized: boolean
}
