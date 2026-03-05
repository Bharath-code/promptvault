document.addEventListener("DOMContentLoaded", (event) => {
    // Register GSAP plugins
    gsap.registerPlugin(ScrollTrigger);

    // Initial Reveal Animations (Hero Section)
    const revealElements = document.querySelectorAll(".gs-reveal");
    
    gsap.fromTo(revealElements, 
        { 
            y: 30, 
            opacity: 0,
            filter: "blur(10px)"
        }, 
        { 
            y: 0, 
            opacity: 1, 
            filter: "blur(0px)",
            duration: 1.2, 
            stagger: 0.15,
            ease: "power3.out",
            delay: 0.2
        }
    );

    // Scroll Fade Up Animations
    const fadeUpElements = document.querySelectorAll(".gs-fade-up");
    
    fadeUpElements.forEach((el) => {
        gsap.fromTo(el,
            {
                y: 50,
                opacity: 0
            },
            {
                scrollTrigger: {
                    trigger: el,
                    start: "top 85%",
                    toggleActions: "play none none reverse"
                },
                y: 0,
                opacity: 1,
                duration: 1,
                ease: "power3.out"
            }
        );
    });

    // Terminal mock typing effect slightly delayed
    const terminalLines = document.querySelectorAll('.terminal-body .line');
    gsap.fromTo(terminalLines,
        {
            opacity: 0,
            x: -10
        },
        {
            scrollTrigger: {
                trigger: ".terminal-mockup",
                start: "top 80%"
            },
            opacity: 1,
            x: 0,
            duration: 0.5,
            stagger: 0.4,
            ease: "power2.out",
            delay: 0.3
        }
    );

    // Mouse movement parallax for glow orbs
    const orbs = document.querySelectorAll('.glow-orb');
    document.addEventListener('mousemove', (e) => {
        const x = e.clientX / window.innerWidth - 0.5;
        const y = e.clientY / window.innerHeight - 0.5;

        gsap.to(orbs[0], {
            x: x * 100,
            y: y * 100,
            duration: 2,
            ease: "power2.out"
        });

        gsap.to(orbs[1], {
            x: -x * 80,
            y: -y * 80,
            duration: 2.5,
            ease: "power2.out"
        });
    });

    // Copy to clipboard functionality
    const copyBtns = document.querySelectorAll('.copy-btn');
    
    copyBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const codeToCopy = btn.parentElement.querySelector('code').innerText;
            
            navigator.clipboard.writeText(codeToCopy).then(() => {
                // Change icon to a checkmark temporarily
                const originalSvg = btn.innerHTML;
                btn.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`;
                btn.classList.add('copied');
                
                // Track copy event with GoatCounter if available
                if (window.goatcounter && window.goatcounter.count) {
                    window.goatcounter.count({
                        path: 'copy-install-cmd',
                        title: 'Copied Install Command',
                        event: true
                    });
                }

                setTimeout(() => {
                    btn.innerHTML = originalSvg;
                    btn.classList.remove('copied');
                }, 2000);
            });
        });
    });
});
