"use client"

import { useQuery } from "react-query"
import { JobPostsTable } from "./table"
import { listMyJobPosts } from "../../queries"
import { JobPost } from "../../types"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation"
import Loading from "../../components/loading"
import Error from "../../components/error"


export default function MyPosts() {
    const router = useRouter()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    if (!cookies.user_token) {
        router.push("/portal/login")
    }
    const query = useQuery<JobPost[], Error>({
        queryKey: ['jobposts', cookies.user_token],
        queryFn: () => listMyJobPosts({ authToken: cookies.user_token }),
    })
    if (query.status == "loading") { return <Loading /> }
    if (query.status == "error") {
        return <Error msg={query.error.message} />
    }
    const jobPosts = Array.from(query.data ? query.data : [])
    return (
        <JobPostsTable jobPosts={jobPosts} />
    )
}
