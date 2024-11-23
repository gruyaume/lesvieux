"use client"

import { useQuery } from "react-query"
import { listAdminAccounts } from "../../queries"
import { UserEntry } from "../../types"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation"
import { UsersTable } from "./table"
import Loading from "../../components/loading"
import Error from "../../components/error"
import { useAuth } from "../auth/authContext"


export default function AdminUsers() {
    const auth = useAuth();
    const router = useRouter()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    if (!cookies.user_token) {
        router.push("/admin_portal/login")
    }
    const adminUsersQuery = useQuery<UserEntry[], Error>({
        queryKey: ['admin_users', cookies.user_token],
        queryFn: () => listAdminAccounts({ authToken: cookies.user_token }),
        retry: (failureCount, error): boolean => {
            if (error.message.includes("401")) {
                return false
            }
            return true
        },
    })

    // if (adminUsersQuery.status == "loading") { return <Loading /> }
    if (adminUsersQuery.status == "error") {
        if (adminUsersQuery.error.message.includes("401")) {
            removeCookie("user_token")
        }
        return <Error msg={adminUsersQuery.error.message} />
    }
    const adminUsers = Array.from(adminUsersQuery.data ? adminUsersQuery.data : [])
    return <UsersTable adminUsers={adminUsers} />
}