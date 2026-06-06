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
                return "border-accent text-accent-foreground bg-accent/10";
            case "success":
                return "border-success-border text-success bg-success-muted";
            case "warning":
                return "border-tertiary/50 text-tertiary bg-tertiary/10";
            case "error":
                return "border-destructive/50 text-destructive bg-destructive/10";
            default:
                return "border-tertiary/50 text-tertiary bg-tertiary/10";
        }
    };
</script>

<div class="border-1 p-4 rounded-lg mb-6 {getAlertTypeStyle(type)}">
    {message}
</div>
