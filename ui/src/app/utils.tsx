
// Utility function to format the date into "YYYY-MM-DD" format
export const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toISOString().split('T')[0];
};

export const HTTPStatus = (code: number): string => {
    const map: { [key: number]: string } = {
        400: "Bad Request",
        401: "Unauthorized",
        403: "Forbidden",
        404: "Not Found",
        409: "Conflict",
        500: "Internal Server Error",
    }
    if (!(code in map)) {
        throw new Error("code not recognized: " + code)
    }
    return map[code]
}

export const passwordIsValid = (pw: string) => {
    if (pw.length < 8) return false

    const result = {
        hasCapital: false,
        hasLowercase: false,
        hasSymbol: false,
        hasNumber: false
    };

    if (/[A-Z]/.test(pw)) {
        result.hasCapital = true;
    }
    if (/[a-z]/.test(pw)) {
        result.hasLowercase = true;
    }
    if (/[0-9]/.test(pw)) {
        result.hasNumber = true;
    }
    if (/[^A-Za-z0-9]/.test(pw)) {
        result.hasSymbol = true;
    }

    if (result.hasCapital && result.hasLowercase && (result.hasSymbol || result.hasNumber)) {
        return true
    }
    return false
}
