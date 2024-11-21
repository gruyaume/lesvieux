"use client"

import { useQuery } from "react-query"
import { listEmployers } from "../../queries"
import { EmployerEntry } from "../../types"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation"
import { EmployersTable } from "./table"
import Loading from "../../components/loading"
import Error from "../../components/error"
import { useAuth } from "../auth/authContext"


export default function Employers() {
    const auth = useAuth();
    const router = useRouter()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    if (!cookies.user_token) {
        router.push("/admin_portal/login")
    }

    const query = useQuery<EmployerEntry[], Error>({
        queryKey: ['employers', cookies.user_token],
        queryFn: () => listEmployers({ authToken: cookies.user_token }),
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
    const employers = Array.from(query.data ? query.data : [])
    return <EmployersTable employers={employers} />
}