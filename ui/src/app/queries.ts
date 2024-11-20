import { BlogPost, UserEntry } from "./types"
import { HTTPStatus } from "./portal/utils"

export type RequiredBlogPostParams = {
    id: string
    authToken: string
}

export async function getStatus() {
    const response = await fetch("/status")
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}


export async function ListPublicBlogPosts(): Promise<BlogPost[]> {
    // Step 1: Fetch the list of blog post IDs
    const response = await fetch("/api/v1/published_posts", {
        method: 'GET',
    });
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}`);
    }

    // The response should look like:
    // {
    //  "result": [1, 2, 3]
    //  }
    const respData = await response.json();
    const ids: number[] = respData.result;

    // Step 2: Fetch details for each blog post
    const promises = ids.map(async (id: number) => {
        const postResponse = await fetch(`/api/v1/published_posts/${id}`, {
            method: 'GET',
        });
        if (!postResponse.ok) {
            throw new Error(`${postResponse.status}: ${HTTPStatus(postResponse.status)}`);
        }

        // The response should look like:
        // {
        //     "result": {
        //         "id": 1,
        //         "title": "abcd",
        //         "content": "abcd",
        //         "created_at": "2024-09-20T18:33:41-04:00",
        //         "author": gruyaume
        //     }
        // }
        const postRespData = await postResponse.json();

        return { ...postRespData.result };
    });

    return Promise.all(promises);
}



export async function listBlogPosts(params: { authToken: string }): Promise<BlogPost[]> {
    // Step 1: Fetch the list of blog post IDs
    const response = await fetch("/api/v1/posts", {
        method: 'GET',
        headers: {
            'Authorization': "Bearer " + params.authToken
        },
    });
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}`);
    }

    // The response should look like:
    // {
    //  "result": [1, 2, 3]
    //  }
    const respData = await response.json();
    const ids: number[] = respData.result;

    // Step 2: Fetch details for each blog post
    const promises = ids.map(async (id: number) => {
        const postResponse = await fetch(`/api/v1/posts/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': "Bearer " + params.authToken
            },
        });
        if (!postResponse.ok) {
            throw new Error(`${postResponse.status}: ${HTTPStatus(postResponse.status)}`);
        }

        // The response should look like:
        // {
        //     "result": {
        //         "id": 1,
        //         "title": "abcd",
        //         "content": "abcd",
        //         "created_at": "2024-09-20T18:33:41-04:00",
        //         "author": gruyaume
        //     }
        // }
        const postRespData = await postResponse.json();

        return { ...postRespData.result };
    });

    return Promise.all(promises);
}



export async function listMyBlogPosts(params: { authToken: string }): Promise<BlogPost[]> {
    // Step 1: Fetch the list of blog post IDs
    const response = await fetch("/api/v1/me/posts", {
        method: 'GET',
        headers: {
            'Authorization': "Bearer " + params.authToken
        },
    });
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}`);
    }

    // The response should look like:
    // {
    //  "result": [1, 2, 3]
    //  }
    const respData = await response.json();
    const ids: number[] = respData.result;

    // Step 2: Fetch details for each blog post
    const promises = ids.map(async (id: number) => {
        const postResponse = await fetch(`/api/v1/me/posts/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': "Bearer " + params.authToken
            },
        });
        if (!postResponse.ok) {
            throw new Error(`${postResponse.status}: ${HTTPStatus(postResponse.status)}`);
        }

        // The response should look like:
        // {
        //     "result": {
        //         "id": 1,
        //         "title": "abcd",
        //         "content": "abcd",
        //         "created_at": "2024-09-20T18:33:41-04:00",
        //         "author": gruyaume
        //     }
        // }
        const postRespData = await postResponse.json();

        return { ...postRespData.result };
    });

    return Promise.all(promises);
}

export async function getBlogPost(params: RequiredBlogPostParams): Promise<BlogPost> {
    const postResponse = await fetch(`/api/v1/posts/${params.id}`, {
        method: 'GET',
        headers: {
            'Authorization': "Bearer " + params.authToken
        },
    });
    if (!postResponse.ok) {
        throw new Error(`${postResponse.status}: ${HTTPStatus(postResponse.status)}`);
    }

    // The response should look like:
    // {
    //     "result": {
    //         "id": 1,
    //         "title": "abcd",
    //         "content": "abcd",
    //         "created_at": "2024-09-20T18:33:41-04:00",
    //         "author": gruyaume
    //     }
    // }
    const postRespData = await postResponse.json();

    return { ...postRespData.result };
}


export async function createBlogPost(params: { authToken: string }) {
    const response = await fetch("/api/v1/me/posts", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer " + params.authToken
        },
        body: JSON.stringify({ "status": "draft" })
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function updateMyBlogPost(params: RequiredBlogPostParams & { title: string, content: string, status: string }) {
    if (!params.title) {
        throw new Error('title not provided')
    }
    if (!params.content) {
        throw new Error('content not provided')
    }
    const response = await fetch("/api/v1/me/posts/" + params.id, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer " + params.authToken
        },
        body: JSON.stringify({ "title": params.title, "content": params.content, "status": params.status })
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function deleteBlogPost(params: RequiredBlogPostParams) {
    const response = await fetch("/api/v1/posts/" + params.id, {
        method: 'DELETE',
        headers: {
            'Authorization': "Bearer " + params.authToken
        }
    })
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${response.statusText}`)
    }
    return
}

export async function deleteMyBlogPost(params: RequiredBlogPostParams) {
    const response = await fetch("/api/v1/me/posts/" + params.id, {
        method: 'DELETE',
        headers: {
            'Authorization': "Bearer " + params.authToken
        }
    })
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${response.statusText}`)
    }
    return
}

export async function login(userForm: { username: string, password: string }) {
    const response = await fetch("/api/v1/login", {
        method: "POST",

        body: JSON.stringify({ "username": userForm.username, "password": userForm.password })
    })
    // The response should look like:
    // {"result":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJncnV5YXVtZSIsInBlcm1pc3Npb25zIjoxLCJleHAiOjE3MjY5NTY0NjV9.oXnHA7YD8Lm-L1iIYAsqhzPXUGTMgOquCkH5XaGERHs"}}
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function changeMyPassword(changePasswordForm: { authToken: string, password: string }) {
    const response = await fetch("/api/v1/me/change_password", {
        method: "POST",
        headers: {
            'Authorization': 'Bearer ' + changePasswordForm.authToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ "password": changePasswordForm.password })
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function changePassword(changePasswordForm: { authToken: string, id: string, password: string }) {
    const response = await fetch("/api/v1/accounts/" + changePasswordForm.id + "/change_password", {
        method: "POST",
        headers: {
            'Authorization': 'Bearer ' + changePasswordForm.authToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ "password": changePasswordForm.password })
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function listUsers(params: { authToken: string }): Promise<UserEntry[]> {
    const response = await fetch("/api/v1/accounts", {
        headers: { "Authorization": "Bearer " + params.authToken }
    })
    // The response should look like:
    // {"result":[{"id":1,"username":"gruyaume","role":1}]}
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function deleteUser(params: { authToken: string, id: string }) {
    const response = await fetch("/api/v1/accounts/" + params.id, {
        method: 'DELETE',
        headers: {
            'Authorization': "Bearer " + params.authToken
        }
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function postFirstUser(userForm: { username: string, password: string }) {
    const response = await fetch("/api/v1/accounts", {
        method: "POST",
        body: JSON.stringify({ "username": userForm.username, "password": userForm.password }),
        headers: {
            'Content-Type': 'application/json'
        }
    })

    // The response should look like:
    // {"result":{"id":1}}
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function postUser(userForm: { authToken: string, username: string, password: string }) {
    const response = await fetch("/api/v1/accounts", {
        method: "POST",
        body: JSON.stringify({
            "username": userForm.username, "password": userForm.password
        }),
        headers: {
            'Authorization': "Bearer " + userForm.authToken
        }
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function isLoggedIn(authToken: string): Promise<boolean> {
    const response = await fetch("/api/v1/me", {
        method: 'GET',
        headers: {
            'Authorization': "Bearer " + authToken
        }
    });

    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return true
}