export default {
    inserted: function (el, binding) {
        el.onclick = function () {
            let total;
            if (binding.value == 0) {
                total = 0;
            } else {

                total = document.getElementById(`anchor-${binding.value}`).offsetTop;
            }
            let distance = document.documentElement.scrollTop || document.body.scrollTop;
            let step = total / 50;
            if (total > distance) {
                (function smoothDown() {
                    if (distance < total) {
                        distance += step;
                        document.documentElement.scrollTop = distance;
                        setTimeout(smoothDown, 5);
                    } else {
                        document.documentElement.scrollTop = total;
                    }
                })();
            } else {
                let newTotal = distance - total;
                step = newTotal / 50;
                (function smoothUp() {
                    if (distance > total) {
                        distance -= step;
                        document.documentElement.scrollTop = distance;
                        setTimeout(smoothUp, 5);
                    } else {
                        document.documentElement.scrollTop = total;
                    }
                })();
            }

        }
    }
}