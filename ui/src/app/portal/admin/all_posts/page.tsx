"use client"

import { useQuery } from "react-query"
import { BlogPostsTable } from "./table"
import { listBlogPosts } from "../../../queries"
import { BlogPost } from "../../../types"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation"
import Loading from "../../../components/loading"
import Error from "../../../components/error"
import { useAuth } from "../../auth/authContext"


export default function BlogPosts() {
    const auth = useAuth();
    const router = useRouter()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    if (!cookies.user_token) {
        router.push("/portal/login")
    }
    if (auth.user && auth.user.role != 1) {
        router.push("/portal/my_posts")
    }
    const query = useQuery<BlogPost[], Error>({
        queryKey: ['blogposts', cookies.user_token],
        queryFn: () => listBlogPosts({ authToken: cookies.user_token }),
    })
    if (query.status == "loading") { return <Loading /> }
    if (query.status == "error") {
        return <Error msg={query.error.message} />
    }
    const blogPosts = Array.from(query.data ? query.data : [])
    return (
        <BlogPostsTable blogPosts={blogPosts} />
    )
}
