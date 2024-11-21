import { JobPost, UserEntry } from "./types"
import { HTTPStatus } from "./utils"

export type RequiredJobPostParams = {
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


export async function ListPublicJobPosts(): Promise<JobPost[]> {
    // Step 1: Fetch the list of job post IDs
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

    // Step 2: Fetch details for each job post
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



export async function listJobPosts(params: { authToken: string }): Promise<JobPost[]> {
    // Step 1: Fetch the list of job post IDs
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

    // Step 2: Fetch details for each job post
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



export async function listMyJobPosts(params: { authToken: string }): Promise<JobPost[]> {
    // Step 1: Fetch the list of job post IDs
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

    // Step 2: Fetch details for each job post
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

export async function getJobPost(params: RequiredJobPostParams): Promise<JobPost> {
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


export async function createJobPost(params: { authToken: string }) {
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

export async function updateMyJobPost(params: RequiredJobPostParams & { title: string, content: string, status: string }) {
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

export async function deleteJobPost(params: RequiredJobPostParams) {
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

export async function deleteMyJobPost(params: RequiredJobPostParams) {
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

export async function employerLogin(userForm: { email: string, password: string }) {
    const response = await fetch("/api/v1/employers/login", {
        method: "POST",

        body: JSON.stringify({ "email": userForm.email, "password": userForm.password })
    })
    // The response should look like:
    // {"result":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJncnV5YXVtZSIsInBlcm1pc3Npb25zIjoxLCJleHAiOjE3MjY5NTY0NjV9.oXnHA7YD8Lm-L1iIYAsqhzPXUGTMgOquCkH5XaGERHs"}}
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function adminLogin(userForm: { email: string, password: string }) {
    const response = await fetch("/api/v1/admin/login", {
        method: "POST",

        body: JSON.stringify({ "email": userForm.email, "password": userForm.password })
    })
    // The response should look like:
    // {"result":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJncnV5YXVtZSIsInBlcm1pc3Npb25zIjoxLCJleHAiOjE3MjY5NTY0NjV9.oXnHA7YD8Lm-L1iIYAsqhzPXUGTMgOquCkH5XaGERHs"}}
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function changeMyEmployerAccountPassword(changeEmployerAccountPasswordForm: { authToken: string, password: string }) {
    const response = await fetch("/api/v1/employers/accounts/me/change_password", {
        method: "POST",
        headers: {
            'Authorization': 'Bearer ' + changeEmployerAccountPasswordForm.authToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ "password": changeEmployerAccountPasswordForm.password })
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function changeMyAdminAccountPassword(changeAdminAccountPasswordForm: { authToken: string, password: string }) {
    const response = await fetch("/api/v1/admin/accounts/me/change_password", {
        method: "POST",
        headers: {
            'Authorization': 'Bearer ' + changeAdminAccountPasswordForm.authToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ "password": changeAdminAccountPasswordForm.password })
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function changeEmployerAccountPassword(changeEmployerAccountPasswordForm: { authToken: string, id: string, password: string }) {
    const response = await fetch("/api/v1/employers/accounts/" + changeEmployerAccountPasswordForm.id + "/change_password", {
        method: "POST",
        headers: {
            'Authorization': 'Bearer ' + changeEmployerAccountPasswordForm.authToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ "password": changeEmployerAccountPasswordForm.password })
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function changeAdminAccountPassword(changeAdminAccountPasswordForm: { authToken: string, id: string, password: string }) {
    const response = await fetch("/api/v1/admin/accounts/" + changeAdminAccountPasswordForm.id + "/change_password", {
        method: "POST",
        headers: {
            'Authorization': 'Bearer ' + changeAdminAccountPasswordForm.authToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ "password": changeAdminAccountPasswordForm.password })
    })
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function listEmployerAccounts(params: { authToken: string }): Promise<UserEntry[]> {
    const response = await fetch("/api/v1/employers/accounts", {
        headers: { "Authorization": "Bearer " + params.authToken }
    })
    // The response should look like:
    // {"result":[{"id":1,"email":"gruyaume","role":1}]}
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function listAdminAccounts(params: { authToken: string }): Promise<UserEntry[]> {
    const response = await fetch("/api/v1/admin/accounts", {
        headers: { "Authorization": "Bearer " + params.authToken }
    })
    // The response should look like:
    // {"result":[{"id":1,"email":"gruyaume","role":1}]}
    const respData = await response.json()
    if (!response.ok) {
        throw new Error(`${response.status}: ${HTTPStatus(response.status)}. ${respData.error}`)
    }
    return respData.result
}

export async function deleteEmployerAccount(params: { authToken: string, id: string }) {
    const response = await fetch("/api/v1/employers/accounts/" + params.id, {
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

export async function deleteAdminAccount(params: { authToken: string, id: string }) {
    const response = await fetch("/api/v1/admin/accounts/" + params.id, {
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

export async function postFirstAdminUser(userForm: { email: string, password: string }) {
    const response = await fetch("/api/v1/admin/accounts", {
        method: "POST",
        body: JSON.stringify({ "email": userForm.email, "password": userForm.password }),
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

export async function postUser(userForm: { authToken: string, email: string, password: string }) {
    const response = await fetch("/api/v1/employers", {
        method: "POST",
        body: JSON.stringify({
            "email": userForm.email, "password": userForm.password
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
    const response = await fetch("/api/v1/admin/accounts/me", {
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