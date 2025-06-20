# Define the directory to process
$directory = ".\gowebserver\react-vite-poc\public"

# Get all files in the directory
Get-ChildItem -Path $directory -File | ForEach-Object {
    # Store the original file name
    $originalName = $_.Name

    # Remove spaces, open parentheses, and close parentheses
    $newName = $originalName -replace " ", "" -replace "\(", "" -replace "\)", ""

    # Rename the file if the name has changed
    if ($originalName -ne $newName) {
        Rename-Item -Path $_.FullName -NewName $newName
        Write-Host "Renamed '$originalName' to '$newName'"
    }
}
