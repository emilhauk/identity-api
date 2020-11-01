class Site {
    constructor(root) {
        this.root = root
    }

    async init() {
        const isLoggedIn = await this.getLoggedInStatus()
        this.hide(this.root.querySelector('.spinner'))
        if(isLoggedIn) {
            const logout = this.root.querySelector('.logout')
            logout.querySelector('h1').textContent = `Hi ${this.user.name}`
            this.show(logout)
            logout.querySelector('[data-action=logout]').addEventListener('click', () => {
                this.logout()
            })
        } else {
            const login = this.root.querySelector('.login')
            this.show(login)
            const form = login.querySelector('form')
            const requestedUrl = new URLSearchParams(location.search).get('requestedUrl')
            if (requestedUrl) {
                const requestedUrlInput = document.createElement('input')
                requestedUrlInput.name = 'requested_url'
                requestedUrlInput.type = 'hidden'
                requestedUrlInput.value = requestedUrl
                form.appendChild(requestedUrlInput)
            }
        }
    }

    getLoggedInStatus() {
        return fetch('/jwt', { method: 'post' }).then(res => {
            if (!res.ok) {
                return false;
            }
            this.bearer = res.headers.get("Authorization")
            this.user = JSON.parse(atob(this.bearer.split('.')[1]))
            return true
        })
    }

    logout() {
        return fetch('/logout', { method: 'POST' })
            .then(res => {
                if (res.ok) {
                    location.reload()
                }
            })
    }

    show(element) {
        element.classList.remove('hidden')
    }

    hide(element) {
        element.classList.add('hidden')
    }

}

let site
document.addEventListener('readystatechange', _ => {
    if (document.readyState === "complete") {
        site = new Site(document.querySelector('.main')).init()
    }
})
