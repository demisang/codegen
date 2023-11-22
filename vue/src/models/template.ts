import Placeholder from "@/models/placeholder";

export default interface Template {
    id: string
    name: string
    description: string
    target_dir: string
    placeholders: Array<Placeholder>
}
