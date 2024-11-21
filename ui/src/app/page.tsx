"use client";

import { useQuery } from "react-query";
import ReaderNavigation from "./nav";
import { formatDate } from "./utils";
import { useRouter } from "next/navigation";
import { useAuth } from "./employer_portal/auth/authContext";
import { getStatus, ListPublicJobPosts } from "./queries";
import { statusResponseResult } from "./types";
import { remark } from 'remark';
import html from 'remark-html';
import { useState, useEffect } from 'react';

interface JobPost {
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
    const { data: jobPosts, isLoading, isError, error } = useQuery<JobPost[], Error>({
        queryKey: "jobPosts",
        queryFn: ListPublicJobPosts,
    });

    const [processedPosts, setProcessedPosts] = useState<JobPost[]>([]);

    useEffect(() => {
        const processContent = async () => {
            if (jobPosts) {
                const postsWithHtml = await Promise.all(jobPosts.map(async (post) => {
                    const processedContent = await remark()
                        .use(html)
                        .process(post.content);
                    return { ...post, content: processedContent.toString() };
                }));
                setProcessedPosts(postsWithHtml);
            }
        };

        processContent();
    }, [jobPosts]);

    if (!auth.firstUserCreated && (statusQuery.data && !statusQuery.data.initialized)) {
        router.push("/employer_portal/initialize");
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
                        <div className="job-post-content" dangerouslySetInnerHTML={{ __html: post.content }} />
                        <hr />
                    </div>
                ))}
            </section>
        </ReaderNavigation>
    );
}
