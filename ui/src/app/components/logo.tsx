"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";

export default function Logo() {
    const router = useRouter();

    const handleLogoClick = () => {
        router.push("/");
    };

    return (
        <div
            onClick={handleLogoClick}
            style={{
                display: "flex",
                alignItems: "center",
                cursor: "pointer",
            }}
        >
            <Image
                alt="LesVieux"
                src="https://www.svgrepo.com/show/523579/notebook-bookmark.svg"
                width="32"
                height="32"
            />
            <span
                style={{
                    marginLeft: "8px",
                    fontSize: "1.5rem",
                    fontWeight: "bold",
                    color: "inherit",
                }}
            >
                LesVieux
            </span>
        </div>
    );
}
