
document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById('searchForm')
    const results = document.getElementById('resultsBox')

    form.addEventListener('submit', async (e) => {
        e.preventDefault(); //Dont Navigate on form Submission

        const data = new URLSearchParams(new FormData(form)); //Channel/term/user

        // 'animate' Searching... 
        let dots = 0;
        results.textContent ="Searching..."
        const interval = setInterval(() => {
           dots = (dots + 1) % 4;
           results.textContent = "Searching" + ".".repeat(dots);
        }, 500);


        const resp = await fetch('/api/search', {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded'},
            body: data
        });

        const text = await resp.text();
        clearInterval(interval);
        results.textContent = text;
    });
});
