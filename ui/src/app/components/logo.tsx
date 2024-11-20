"use client"

import Image from "next/image";

export default function Logo() {
    return (
        <div style={{ display: "flex", alignItems: "center" }}>
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
    )
}
