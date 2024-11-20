"use client";

import { useQuery } from "react-query";
import ReaderNavigation from "./nav";
import { formatDate } from "./utils";
import { useRouter } from "next/navigation";
import { useAuth } from "./portal/auth/authContext";
import { getStatus, ListPublicBlogPosts } from "./queries";
import { statusResponseResult } from "./types";
import { remark } from 'remark';
import html from 'remark-html';
import { useState, useEffect } from 'react';

interface BlogPost {
    id: number;
    title: string;
    content: string;
    created_at: string;  // RFC3339 date format
    author: string;
}

export default function FrontPage() {
    const router = useRouter();
    const auth = useAuth();
    const statusQuery = useQuery<statusResponseResult, Error>({
        queryKey: "status",
        queryFn: () => getStatus()
    });
    const { data: blogPosts, isLoading, isError, error } = useQuery<BlogPost[], Error>({
        queryKey: "blogPosts",
        queryFn: ListPublicBlogPosts,
    });

    const [processedPosts, setProcessedPosts] = useState<BlogPost[]>([]);

    useEffect(() => {
        const processContent = async () => {
            if (blogPosts) {
                const postsWithHtml = await Promise.all(blogPosts.map(async (post) => {
                    const processedContent = await remark()
                        .use(html)
                        .process(post.content);
                    return { ...post, content: processedContent.toString() };
                }));
                setProcessedPosts(postsWithHtml);
            }
        };

        processContent();
    }, [blogPosts]);

    if (!auth.firstUserCreated && (statusQuery.data && !statusQuery.data.initialized)) {
        router.push("/portal/initialize");
    }

    return (
        <ReaderNavigation>
            <section className="p-section">
                {isLoading && <div>Loading...</div>}

                {isError && <div>Error: {error?.message}</div>}

                {!isLoading && !isError && processedPosts?.map((post) => (
                    <div key={post.id} className="p-section--shallow">
                        <h2>{post.title}</h2>
                        <h5>By: {post.author}</h5>
                        <h6>{formatDate(post.created_at)}</h6>
                        <div className="blog-post-content" dangerouslySetInnerHTML={{ __html: post.content }} />
                        <hr />
                    </div>
                ))}
            </section>
        </ReaderNavigation>
    );
}
