"use client"

import { useQuery } from "react-query"
import { listUsers } from "../../../queries"
import { UserEntry } from "../../../types"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation"
import { UsersTable } from "./table"
import Loading from "../../../components/loading"
import Error from "../../../components/error"
import { useAuth } from "../../auth/authContext"


export default function Users() {
    const auth = useAuth();
    const router = useRouter()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    if (!cookies.user_token) {
        router.push("/portal/login")
    }
    if (auth.user && auth.user.role != 1) {
        router.push("/portal/my_posts")
    }
    const query = useQuery<UserEntry[], Error>({
        queryKey: ['users', cookies.user_token],
        queryFn: () => listUsers({ authToken: cookies.user_token }),
        retry: (failureCount, error): boolean => {
            if (error.message.includes("401")) {
                return false
            }
            return true
        },
    })
    if (query.status == "loading") { return <Loading /> }
    if (query.status == "error") {
        if (query.error.message.includes("401")) {
            removeCookie("user_token")
        }
        return <Error msg={query.error.message} />
    }
    const users = Array.from(query.data ? query.data : [])
    return <UsersTable users={users} />
}