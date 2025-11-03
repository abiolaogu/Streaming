# GitHub CLI Authentication Guide

## Quick Authentication Steps

GitHub CLI is already installed on your system! Here's how to authenticate:

### Method 1: Web Browser (Recommended - Easiest)

1. **Open Terminal** and run:
   ```bash
   gh auth login
   ```

2. **Follow the prompts**:
   - When asked "What account do you want to log into?", choose:
     ```
     GitHub.com
     ```
     (Press Enter or type `1`)

   - When asked "What is your preferred protocol?", choose:
     ```
     HTTPS
     ```
     (Press Enter or type `1`)

   - When asked "How would you like to authenticate?", choose:
     ```
     Login with a web browser
     ```
     (Press Enter or type `1`)

   - When asked "Paste your authentication token:", the terminal will display:
     ```
     ! First copy your one-time code: XXXX-XXXX
     ```
     - **Copy the code** (it will look like: `ABCD-1234`)
     - Press Enter
     - Your **web browser will automatically open** to GitHub
     - If it doesn't, go to: https://github.com/login/device
     - **Paste the code** into the GitHub page
     - Click "Authorize"
     - You'll see "Authentication complete. Press Enter to continue"
     - Go back to Terminal and **press Enter**

3. **Verify Authentication**:
   ```bash
   gh auth status
   ```
   You should see: `✓ Logged in to github.com as YOUR_USERNAME`

### Method 2: Personal Access Token (If Browser Doesn't Work)

1. **Create a Personal Access Token on GitHub**:
   - Go to: https://github.com/settings/tokens
   - Click "Generate new token" → "Generate new token (classic)"
   - Name it: "GitHub CLI"
   - Select scopes:
     - ✅ `repo` (Full control of private repositories)
     - ✅ `read:org` (if working with organizations)
   - Click "Generate token"
   - **Copy the token immediately** (you won't see it again!)

2. **Authenticate with Token**:
   ```bash
   gh auth login
   ```
   - Choose: `GitHub.com`
   - Choose: `HTTPS`
   - Choose: `Paste an authentication token`
   - Paste your token and press Enter

3. **Verify**:
   ```bash
   gh auth status
   ```

## After Authentication

Once authenticated, you can create GitHub issues:

```bash
# Test with dry run first
python3 create_github_issues.py --dry-run --range 1-3

# Create all issues
python3 create_github_issues.py
```

## Troubleshooting

### "Command not found: gh"
If you see this, GitHub CLI might not be in your PATH. Try:
```bash
# Add to PATH (for zsh - default on macOS)
echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### "Permission denied"
Make sure you have write access to the repository:
```bash
gh repo view abiolaogu/https-github.com-abiolaogu-Video-Streaming_AI-Studio
```

### "Authentication failed"
- Try logging out and back in:
  ```bash
  gh auth logout
  gh auth login
  ```
- Make sure you have a GitHub account
- Check your internet connection

## Need Help?

If you're stuck, you can also authenticate via:
- GitHub Desktop (if installed)
- Manual token creation on GitHub.com

Let me know if you need help with any step!

