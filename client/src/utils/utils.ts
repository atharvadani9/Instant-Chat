export const getUserID = () => {
    return localStorage.getItem("userID");
}

export const getAvatarColor = (username: string) => {
    const colors = [
        '#7C4DFF', // Light Blue
        '#00BCD4', // Light Green
        '#FF9800', // Light Orange
        '#E91E63', // Light Pink
        '#9C27B0', // Light Purple
        '#009688'  // Light Teal
    ];

    // Generate a hash from the username for consistency
    let hash = 0;
    for (let i = 0; i < username.length; i++) {
        hash = username.charCodeAt(i) + ((hash << 5) - hash);
    }

    // Select color based on hash
    const colorIndex = Math.abs(hash) % colors.length;
    return colors[colorIndex];
};