<script lang="ts">
    type Props = {
        type?: AlertType;
        message: string;
    };
    let { type = "info", message }: Props = $props();

    type AlertType = "info" | "error" | "warning" | "success";
    const getAlertTypeStyle = (type: AlertType): string => {
        switch (type) {
            case "info":
                return "border-blue-500 text-blue-800 bg-blue-50";
            case "success":
                return "border-green-500 text-green-800 bg-green-50";
            case "warning":
                return "border-amber-500 text-amber-800 bg-amber-50";
            case "error":
                return "border-red-500 text-red-800 bg-red-50";
            default:
                return "border-amber-500 text-amber-800 bg-amber-50";
        }
    };
</script>

<div class="border-1 p-4 rounded-lg mb-6 {getAlertTypeStyle(type)}">
    {message}
</div>
