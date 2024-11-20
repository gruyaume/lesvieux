"use client"

import { useQuery } from "react-query"
import { JobPostsTable } from "./table"
import { listJobPosts } from "../../../queries"
import { JobPost } from "../../../types"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation"
import Loading from "../../../components/loading"
import Error from "../../../components/error"
import { useAuth } from "../../auth/authContext"


export default function JobPosts() {
    const auth = useAuth();
    const router = useRouter()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    if (!cookies.user_token) {
        router.push("/portal/login")
    }
    if (auth.user && auth.user.role != 1) {
        router.push("/portal/my_posts")
    }
    const query = useQuery<JobPost[], Error>({
        queryKey: ['jobposts', cookies.user_token],
        queryFn: () => listJobPosts({ authToken: cookies.user_token }),
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
